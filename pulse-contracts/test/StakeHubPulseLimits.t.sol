// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity 0.8.17;

import "../contracts/StakeHub.sol";

interface Vm {
    function prank(address sender) external;
    function expectRevert(bytes4 revertData) external;
}

contract StakeHubLimitHarness is StakeHub {
    function checkSelfStakeLimit(uint256 currentStake, uint256 additionalStake) external pure {
        _checkSelfDelegationLimit(currentStake, additionalStake);
    }
}

contract StakeHubPulseLimitsTest {
    Vm private constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    function testPulseStakeParametersInitialize() external {
        StakeHub hub = new StakeHub();
        vm.prank(block.coinbase);
        hub.initialize();

        require(hub.minSelfDelegationBNB() == 200_000 ether, "wrong PULSE minimum");
        require(hub.MAX_SELF_DELEGATION_PULSE() == 10_000_000 ether, "wrong PULSE maximum");
        require(hub.maxElectedValidators() == 45, "elected-validator capacity changed");
    }

    function testSelfStakeCapAllowsExactMaximum() external {
        StakeHubLimitHarness hub = new StakeHubLimitHarness();
        hub.checkSelfStakeLimit(9_000_000 ether, 1_000_000 ether);
    }

    function testSelfStakeCapRejectsAboveMaximum() external {
        StakeHubLimitHarness hub = new StakeHubLimitHarness();
        vm.expectRevert(StakeHub.SelfDelegationTooMuch.selector);
        hub.checkSelfStakeLimit(9_000_000 ether, 1_000_000 ether + 1);
    }

    function testGovernanceMinimumAcceptsConfiguredBounds() external {
        StakeHub hub = new StakeHub();
        vm.prank(block.coinbase);
        hub.initialize();

        vm.prank(address(0x1007));
        hub.updateParam("minSelfDelegationBNB", abi.encode(uint256(200_000 ether)));
        require(hub.minSelfDelegationBNB() == 200_000 ether, "governance lower bound rejected");

        vm.prank(address(0x1007));
        hub.updateParam("minSelfDelegationBNB", abi.encode(uint256(10_000_000 ether)));
        require(hub.minSelfDelegationBNB() == 10_000_000 ether, "governance upper bound rejected");
    }

    function testGovernanceMinimumRejectsValuesOutsideBounds() external {
        StakeHub hub = new StakeHub();
        vm.prank(block.coinbase);
        hub.initialize();

        vm.expectRevert(StakeHub.InvalidValue.selector);
        vm.prank(address(0x1007));
        hub.updateParam("minSelfDelegationBNB", abi.encode(uint256(199_999 ether)));

        vm.expectRevert(StakeHub.InvalidValue.selector);
        vm.prank(address(0x1007));
        hub.updateParam("minSelfDelegationBNB", abi.encode(uint256(10_000_001 ether)));
    }
}

package pulse

import _ "embed"

// StakeHubContract is the PULSE-specific StakeHub runtime bytecode generated
// from pulse-contracts/contracts/StakeHub.sol by the pinned Foundry build.
//go:embed StakeHubContract
var StakeHubContract string

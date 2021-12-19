package example

// ==================================================
// ============ Example for build info ==============
// ==================================================

// These infos might inject in the build time
// In order to provide build info, add the following option to bootstrap file
// golib.BuildInfoOpt(Version, CommitHash, BuildTime)
var (
	Version    = "1.0"
	CommitHash = "49d52932"
	BuildTime  = "2021/11/30 20:08:12"
)

package constants

var (
	appName       string = "PuGo"
	appVersion    string = "dev"
	appGithubLink        = "https://github.com/fuxiaohei/pugo"
)

// AppName returns app name
func AppName() string {
	return appName
}

// AppVersion returns app version
func AppVersion() string {
	return appVersion
}

// SetAppVersion sets app version
func SetAppVersion(version string) {
	appVersion = version
}

// AppGithubLink returns app github link
func AppGithubLink() string {
	return appGithubLink
}

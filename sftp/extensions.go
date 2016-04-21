package sftp

// SFTP Protocol Extensions
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-13#section-4.2
type ExtensionPair struct {
	Name string
	Data string
}

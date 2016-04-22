package sftp

// SFTP Protocol Extensions
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-13#section-4.2
type ExtensionPair struct {
	Name string
	Data string
}

func UnmarshalExtensionPair(b []byte) (ExtensionPair, []byte, error) {
	var extension ExtensionPair
	var err error

	extension.Name, b, err = UnmarshalStringSafe(b)

	if err != nil {
		return extension, b, err
	}

	extension.Data, b, err = UnmarshalStringSafe(b)

	if err != nil {
		return extension, b, err
	}

	return extension, b, nil
}

package helper

const TagMaxLength = 128

func TruncateTag(tag string) string {
	if len(tag) > TagMaxLength {
		return tag[:TagMaxLength]
	}
	return tag
}

package sender

type (
	Sender interface {
		Send(repoName string, materials Materials) error
	}

	Material struct {
		Title       string
		LinkURL     string
		Reviewers   []string
		IsMergeable bool
	}

	Materials []Material
)

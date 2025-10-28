package ports

// InboundHandler est le port d'entrée (contrat) implémenté par chaque handler CLI.
type InboundHandler interface {
	Kind() string
	Handle(yamlPayload []byte) error
}

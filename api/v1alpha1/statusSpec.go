package v1alpha1

type HttpStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Code int    `json:"code,omitempty"`
	Body string `json:"body,omitempty"`
}

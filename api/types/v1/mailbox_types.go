package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MailboxSpec defines the desired state of Mailbox
type MailboxSpec struct {
	EmailAddress string `json:"emailAddress"`
	PasswordHash string `json:"passwordHash,omitempty"`
	IsForward    bool   `json:"isForward,omitempty"`
	Enabled      bool   `json:"enabled,omitempty"`
}

// MailboxStatus defines the observed state of Mailbox
type MailboxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
}

// +kubebuilder:object:root=true

// Mailbox is the Schema for the mailboxes API
type Mailbox struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MailboxSpec   `json:"spec,omitempty"`
	Status MailboxStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MailboxList contains a list of Mailbox
type MailboxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mailbox `json:"items"`
}

// func init() {
// 	SchemeBuilder.Register(&Mailbox{}, &MailboxList{})
// }

package wallet

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Wallet ...
type Wallet struct {
	ID           bson.ObjectId `bson:"_id"`
	WalletAmount WalletAmount  `bson:"wallet_amount"`
	Transaction  *Transaction
	Transactions []*Transaction `bson:"transactions"`
	CreatedAt    time.Time      `bson:"created_at"`
}

// Transactions ...
type Transactions struct {
	Transactions      []*Transaction `bson:"transactions"`
	TransactionAmount int32          `bson:"transaction_amount"`
}

// Transaction ...
type Transaction struct {
	Type            TransactionType     `bson:"transaction_type"`
	CoinType        TransactionCoinType `bson:"coin_type"`
	Status          Status              `bson:"status"`
	WalletAmount    WalletAmount        `bson:"wallet_amount"`
	TransactionTime time.Time           `bson:"transaction_at"`
}

// TransactionType ...
type TransactionType string

const (
	// TransactionTypeUnknown ...
	TransactionTypeUnknown TransactionType = "unknown"
	// TransactionTypeInvitation ...
	TransactionTypeInvitation TransactionType = "invitation"
	// TransactionTypeShare ...
	TransactionTypeShare TransactionType = "share"
	// TransactionTypeApplyJob ...
	TransactionTypeApplyJob TransactionType = "apply_job"
	// TransactionTypeUserRegistation ...
	TransactionTypeUserRegistation TransactionType = "user_registration"
	// TransactionTypeCompanyRegistation ...
	TransactionTypeCompanyRegistation TransactionType = "company_registation"
	// TransactionTypeBecomeCandidate ...
	TransactionTypeBecomeCandidate TransactionType = "become_candidate"
	// TransactionTypeCreatePost ...
	TransactionTypeCreatePost TransactionType = "create_post"
	// TransactionJobShare ...
	TransactionJobShare TransactionType = "job_share"
)

// TransactionCoinType ...
type TransactionCoinType string

const (
	// TransactionCoinTypeAll ...
	TransactionCoinTypeAll TransactionCoinType = "all"

	// TransactionCoinTypeGold ...
	TransactionCoinTypeGold TransactionCoinType = "gold"

	// TransactionCoinTypeSilver ...
	TransactionCoinTypeSilver TransactionCoinType = "silver"
)

// Status ...
type Status string

const (
	// WalletStatusDone ...
	WalletStatusDone Status = "done"

	// WalletStatusPending ...
	WalletStatusPending Status = "pending"

	// WalletStatusRejected ...
	WalletStatusRejected Status = "rejected"
)

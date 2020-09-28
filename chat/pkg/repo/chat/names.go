package chat

const CONVERSATIONS_COLLECTION_NAME = "conversations"
const MESSAGES_COLLECTION_NAME = "messages"
const LABELS_COLLECTION_NAME = "labels"
const REPLIES_COLLECTIONS_NAME = "replies"
const REPORTS_COLLECTIONS_NAME = "reports"
const FS_NAME = "chat_fs"

const CONVERSATIONS_INDEX_NAME = "participants_idx"
const CONVERSATIONS_INDEX_KEY = "participants.id"

const MESSAGES_INDEX_NAME = "conversations_idx"
const MESSAGES_INDEX_KEY = "conversation_id"

const REPLIES_INDEX_NAME = "replies_idx"
const REPLIES_INDEX_KEY = "user_id"

const LABELS_INDEX_NAME = "labels_idx"
const LABELS_INDEX_KEY = "user_id"

package messagehub

// job def change message from Repository
const TOPIC_JOB_DEFINITION = "job-definition"

// node def change message from Repository
const TOPIC_NODE_DEFINITION = "node-definition"

// job log from NodeManager
const TOPIC_JOB_LOG = "job-log"

// job state change (start job, start step, and end) message from NodeManager
// messageType = *message.ExecuterMsg
const TOPIC_JOB_REPORT = "job-report"

// config change message from Repository
// messageType = *message.ConfigMsg
const TOPIC_CONFIG_CHANGE = "config-change"

// job run request message from JobScheduler
// deprecated: dont use
const TOPIC_JOB_RUN_REQUEST = "job-run-request"

// use for backup
const TOPIC_FREEZE_FILESYSTEM = "freeze-filesystem"

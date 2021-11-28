package messagehub

// job def change message from Repository
const TOPIC_JOB_DEFINITION = "job-definition"

// node def change message from Repository
const TOPIC_NODE_DEFINITION = "node-definition"

// job state change message from NodeManager
const TOPIC_JOB_REPORT = "job-report"

// job run request message from JobScheduler
const TOPIC_JOB_RUN_REQUEST = "job-run-request"

// use for backup
const TOPIC_FREEZE_FILESYSTEM = "freeze-filesystem"

1. Vague dev prompt

    Parse this log and give me the important bits

2. Rewrite of the Prompt

    Extract the following fields from the log line and return only valid JSON.
    If a field is missing, set it to null. Do not invent values.

    Fields:
    - timestamp
    - level
    - thread
    - logger
    - message

    Log line:
    2025-07-14 11:02:33.891 [worker-3] ERROR com.acme.jobs.SyncJob - Failed to sync inventory for warehouse WH-17: connection timed out after 30s



3. Outputs:

    Weak output:

        > Input
        2025-07-14 11:02:33.891 [worker-3] ERROR com.acme.jobs.SyncJob - Failed to sync inventory for warehouse WH-17: connection timed out after 30s

        > Output
        The log shows an ERROR from SyncJob about failing to sync inventory for warehouse WH-17 because of a connection timeout.

    Rewrite:

        > Input
        2025-07-14 11:02:33.891 [worker-3] ERROR com.acme.jobs.SyncJob - Failed to sync inventory for warehouse WH-17: connection timed out after 30s

        > Output
        {
        "timestamp": "2025-07-14 11:02:33.891",
        "level": "ERROR",
        "thread": "worker-3",
        "logger": "com.acme.jobs.SyncJob",
        "message": "Failed to sync inventory for warehouse WH-17: connection timed out after 30s"
        }

    One Sentence:

        The instructional pattern forces a precise schema and null-handling rule, turning free-form prose into reliable, machine-readable JSON.
        It is reusable, and in case of many inputs, all the asnwers would look the same. so it's a lot easier to file them and go through them.
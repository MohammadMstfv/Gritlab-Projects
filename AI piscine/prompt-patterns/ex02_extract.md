1. I chose Log line.


2. here is the template and the prompt:

    Extract these fields from the log line and return JSON:
    {
    "timestamp": ...,
    "level": ...,
    "logger": ...,
    "message": ...,
    "thread": ...
    }
    Log line: {INPUT}


3. Two real examples:

    > Input
    Log line: 2024-09-12 14:23:07.482 [http-nio-8080-exec-5] ERROR c.e.p.s.PaymentService - Payment failed for orderId=ORD-98421: insufficient funds (code=INSUFFICIENT_FUNDS)

    > Output
    {
    "timestamp": "2024-09-12 14:23:07.482",
    "level": "ERROR",
    "logger": "c.e.p.s.PaymentService",
    "message": "Payment failed for orderId=ORD-98421: insufficient funds (code=INSUFFICIENT_FUNDS)",
    "thread": "http-nio-8080-exec-5"
    }

    > Input
    Log line: 2025-03-18T09:41:55.113Z INFO  [main] o.s.b.w.e.tomcat.TomcatWebServer - Tomcat started on port(s): 8080 (http) with context path ''

    > Output
    {
    "timestamp": "2025-03-18T09:41:55.113Z",
    "level": "INFO",
    "logger": "o.s.b.w.e.tomcat.TomcatWebServer",
    "message": "Tomcat started on port(s): 8080 (http) with context path ''",
    "thread": "main"
    }

OPTIONAL: 

    Extract these fields from the log line and return JSON:
    {
    "timestamp": ...,
    "level": ...,
    "logger": ...,
    "message": ...,
    "thread": ...
    }

    Rules:
    - If a field cannot be extracted or is missing from the log line, set its value to null.
    - Do not invent or guess missing values.
    - Always return valid JSON.

    Log line: {INPUT}

    > Input
    Log line: 2025-03-18T09:41:55.113Z WARN Payment declined for user 4421 – card expired

    > Output
    {
    "timestamp": "2025-03-18T09:41:55.113Z",
    "level": "WARN",
    "logger": null,
    "message": "Payment declined for user 4421 – card expired",
    "thread": null
    }
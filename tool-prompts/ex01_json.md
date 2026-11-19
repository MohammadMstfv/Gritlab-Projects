1. Prompt:

    Answer the question below.
    Return ONLY valid JSON, with no extra text or markdown. and use a 0 to 100 number to show how confident are you with the answer.
    Use exactly this shape:
    {"answer": string, "confidence": number}

    Question: {{How many ants are in Brazil?}}

2. Ran:

    {"answer": "No exact count exists. Global estimate is ~20 quadrillion ants; Brazil (large tropical share, Amazon + Cerrado) likely holds roughly 1–3 quadrillion.", "confidence": 25}

    is the Json valid? 
        well it kinda isn't since it used the string part for explanation except just answer.

    so I changed the string to number too, and asked same kinda questions.

    {"answer": 2000000000000000, "confidence": 20}

    Now it's valid.

    How many countries are in the world.
    {"answer": 195, "confidence": 95}

    How many r letters are in the word Strawberry?
    {"answer": 3, "confidence": 100}





3. One Sentence:
    A sentence is optimized for humans reading it. JSON is optimized for code acting on it. When a program is the consumer, you want the data pre-structured so it can be used without interpretation.
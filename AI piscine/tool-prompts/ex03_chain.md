Step 1: User request:

    How much 178cm is in feet

Step 2: Model's JSON tool call (from Exercise 2)

    JSON{"tool": "convert_units", "args": {"value": 178, "from": "cm", "to": "ft"}}

Step 3: Fake tool "result" + request for final answer

    Fake tool "result" :
    JSON{"result": 5.83989501312336, "unit": "ft", "formatted": "5 ft 10.08 in"}


    Prompt to model:

    Using the tool result above, phrase a clear final answer for the user.

Expected model final answer:

    178 cm is approximately 5.84 feet (exactly 5 feet 10.08 inches).

Where this chain could break (one sentence)

    The chain can break if the tool returns malformed/empty JSON, an unexpected unit, or a rate-limited/error response, so the orchestration code must validate the tool payload, handle missing fields, and either retry or return a graceful fallback message before generating the final answer.

Going further (notes)

    Broken/empty JSON handling: Code should catch JSON decode errors, check for required keys (result, etc.), log the failure, and either retry once or refuse with a user-friendly message instead of crashing or hallucinating a number.

    Real wiring: A small Python script can parse the model's JSON, dispatch to a real convert_units implementation (or a mock), then feed the result back into a second LLM call that produces the final natural-language answer. Most provider SDKs already support this via native function/tool calling.

    > Python code

    import json

    # ---------- 1. Fake tool implementations ----------
    def calculator(expression: str):
        # Only allow safe math chars
        allowed = set("0123456789+-*/(). ")
        if not set(expression) <= allowed:
            return {"error": "unsafe expression"}
        return {"result": eval(expression)}

    def search_web(query: str):
        fake_db = {
            "population of france": "France population: ~68,000,000 (2024)",
            "capital of australia": "Canberra",
        }
        return {"result": fake_db.get(query.lower(), f"No result for '{query}'")}

    def run_python(code: str):
        # DANGER: eval of arbitrary code — demo only, sandbox in real life
        local_ns = {}
        try:
            exec(code, {}, local_ns)
            return {"output": {k: v for k, v in local_ns.items() if not k.startswith("__")}}
        except Exception as e:
            return {"error": str(e)}

    def convert_units(value: float, from_unit: str, to_unit: str):
        table = {
            ("km2", "mi2"): 0.386102,
            ("c", "f"): lambda v: v * 9 / 5 + 32,
            ("usd", "jpy"): 150.0,
        }
        key = (from_unit.lower(), to_unit.lower())
        if key not in table:
            return {"error": f"no conversion {from_unit}->{to_unit}"}
        factor = table[key]
        result = factor(value) if callable(factor) else value * factor
        return {"result": round(result, 4)}

    def translate(text: str, target_lang: str):
        fake = {
            ("good morning", "fr"): "Bonjour",
            ("thank you", "ja"): "ありがとう",
        }
        return {"result": fake.get((text.lower(), target_lang.lower()), f"[{target_lang}] {text}")}

    def ask_user(question: str):
        # In a real app this would block for input; here we auto-answer.
        return {"answer": "Use 2024 data."}

    def finish(answer: str):
        return {"__done__": True, "answer": answer}

    # ---------- 2. Registry ----------
    TOOLS = {
        "calculator": calculator,
        "search_web": search_web,
        "run_python": run_python,
        "convert_units": convert_units,
        "translate": translate,
        "ask_user": ask_user,
        "finish": finish,
    }

    def dispatch(tool_call: dict) -> dict:
        """tool_call looks like {"tool": "...", "arguments": {...}}"""
        name = tool_call.get("tool")
        args = tool_call.get("arguments", {})
        if name not in TOOLS:
            return {"error": f"unknown tool '{name}'"}
        try:
            return TOOLS[name](**args)
        except TypeError as e:
            return {"error": f"bad arguments for {name}: {e}"}

    # ---------- 3. The "model" — replace with a real SDK call ----------
    def call_model(messages: list[dict]) -> dict:
        """
        Stand-in for a real LLM call. Returns ONE tool call as JSON.
        In production this is where you'd call OpenAI/Anthropic/etc.
        """
        # Very dumb scripted sequence for demo purposes:
        step = len([m for m in messages if m["role"] == "tool"])
        script = [
            {"tool": "ask_user", "arguments": {"question": "Which year's data?"}},
            {"tool": "search_web", "arguments": {"query": "population of france"}},
            {"tool": "calculator", "arguments": {"expression": "68000000 / 8000000000 * 100"}},
            {"tool": "convert_units", "arguments": {"value": 551695, "from_unit": "km2", "to_unit": "mi2"}},
            {"tool": "run_python", "arguments": {"code": "evens = [n for n in range(10) if n % 2 == 0]"}},
            {"tool": "translate", "arguments": {"text": "good morning", "target_lang": "fr"}},
            {"tool": "finish", "arguments": {"answer": "France is ~68M, ~0.85% of world pop."}},
        ]
        return script[step] if step < len(script) else script[-1]

    # ---------- 4. The agent loop ----------
    def run_agent(user_prompt: str, max_steps: int = 10):
        messages = [{"role": "user", "content": user_prompt}]

        for _ in range(max_steps):
            tool_call = call_model(messages)
            print(f"→ model wants: {json.dumps(tool_call)}")

            result = dispatch(tool_call)
            print(f"← tool result: {json.dumps(result)}\n")

            messages.append({"role": "tool", "tool_call": tool_call, "result": result})

            if result.get("__done__"):
                print("FINAL:", result["answer"])
                return result["answer"]

        print("Hit max steps without finishing.")
        return None

    if __name__ == "__main__":
        run_agent("Give me a quick France stats report.")
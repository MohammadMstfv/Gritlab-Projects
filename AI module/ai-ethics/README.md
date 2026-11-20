# AI Ethics

## Main Takeaway

The main thing I learned from this quest is that building AI responsibly is not just about making a good prompt or making the model give better answers. An AI system can cause problems because of its training data, the way the task is designed, the model itself, or the way the final product is deployed. As a developer, I need to think about what could go wrong, who could be affected, how serious the consequences could be, and what kind of solution is actually appropriate.

I also learned that I cannot use "the model did it" as an excuse when an AI feature causes harm. If I choose a model, build it into a product, decide how users interact with it, and deploy it, I have some responsibility for the risks created by those decisions. Prompting can reduce some problems, but sometimes the real solution needs better data, a different model, product-level safeguards, testing, monitoring, or human oversight.

## What Changed My Thinking

Before this quest, I thought that if an AI produced a harmful or biased answer, improving the prompt would probably be one of the main ways to fix it. After working through the examples, I understand that prompt design is only one layer of an AI system.

For example, if a model has learned a biased pattern from its training data, telling it to "be fair" in the prompt does not necessarily remove the underlying bias. Similarly, if an AI gives dangerously incorrect information, a better prompt might reduce the problem but may not be enough. The system might also need better data, a more reliable model, retrieval from trusted sources, output validation, human review, or restrictions on what the AI is allowed to do.

The Air Canada chatbot case also made the accountability issue more concrete for me. The fact that an AI system generated incorrect information does not mean the organization using it can simply say that the AI was responsible. The way the system was designed and deployed is also part of the responsibility.

## What I Learned

The biggest lesson for me is to look at an AI feature as a complete system rather than just looking at the model's response.

When I see an AI behaving badly, I should ask:

1. Where did the problem come from?
2. Is it related to the data?
3. Is the task or prompt framed badly?
4. Is the model itself unreliable for this task?
5. How is the AI being deployed?
6. Who could be harmed?
7. How serious could the harm be?
8. Would a prompt change actually solve the problem?
9. Do I need a product-level safeguard or human oversight?
10. Who is responsible for preventing the problem?

This way of thinking is something I can apply to future software and AI projects, even when the project is not specifically about ethics.

## Sources

* **OECD AI Principles** — The OECD is an international organization that develops policy standards with governments and provides a framework for trustworthy AI. Its AI Principles cover areas including human rights and fairness, transparency, robustness and safety, and accountability.
* **UNESCO Recommendation on the Ethics of Artificial Intelligence** — This is an international framework adopted by UNESCO Member States that addresses human rights, human dignity, fairness, transparency, human oversight, accountability, safety, and non-discrimination throughout the AI lifecycle.

## Final Reflection

Responsible AI does not mean trying to make an AI system perfect, because AI systems will still make mistakes. It means understanding the risks before and after deployment and designing the system so that predictable mistakes are less likely to cause serious harm.

For me, the most useful mindset is to stop asking only:

> "How can I make the AI give a better answer?"

and also ask:

> "What happens if this answer is wrong, who is affected, and what should I have built around the AI to make that failure safer?"

That is the main idea I want to carry into future AI and software projects.

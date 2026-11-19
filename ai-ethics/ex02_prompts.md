1. Three Problematic Prompts

    1. "Look at these resumes and identify the candidate who seems most naturally suited to leadership."
    > problem
    "Naturally suited" suggests leadership is an inherent personal trait."Leadership" is undefined. "Seems" invites subjective inference. Also Resume information could contain demographic or proxy information that encourages irrelevant judgments. Harm comes from here: the model could make confident judgments about people based on weak evidence.
    > better version
    "Compare these candidates only using the explicitly listed job-related criteria: project management experience, team leadership experience, technical experience, and communication experience. For each candidate, identify the evidence supporting your assessment. Do not infer protected characteristics or personality traits that are not supported by the application materials. If the evidence is insufficient, say so."
    > why is it better
    defined criteria. evidence requirements. reduced subjective inference. limited scope. uncertainty handling. prohibited unsupported assumptions.

    2. "I have these symptoms. Tell me what disease I have and what medication I should take."
    > problem
    asks for diagnosis. asks for treatment. insufficient information. high consequence. encourages false certainty
    > better version
    "Help me understand what kinds of conditions can sometimes cause these symptoms. Clearly distinguish possibilities from diagnoses, identify symptoms that may require urgent medical attention, and suggest what information I should discuss with a qualified healthcare professional. Do not prescribe medication or present a diagnosis as certain."
    > why is it better
    doesn't let it give wrong information confidently. gives more info about it before telling what to do about it...

    3. "Tell me which investment will make me the most money."
    > problem
    assumes future prediction is possible. encourages overconfidence. doesn't define risk. ignores time horizon. ignores user's financial circumstances. asks for a consequential recommendation
    > better version
    "Compare these investments using historical risk, volatility, liquidity, fees, diversification, and time horizon. Clearly distinguish historical information from predictions, identify important uncertainties, and avoid presenting future returns as guaranteed."
    > why is it better
    gives more context, doesn't let it do it on itself, doesn't allow fake confidence


3. The Answers did change yes. but prompt alone doesn't fix the problem.


AI review

    You have done excellent work identifying the core ethical problems in each prompt and constructing thoughtful rewrites. Your reasoning is sound: you correctly identified issues of subjective inference, false certainty, and ungrounded prediction. However, there are three interconnected ethical dimensions you have overlooked. These are not minor gaps—they are central to what makes AI interactions ethically safe.

    1. The Illusion of Sufficiency: You Have Not Addressed the Limits of the AI Itself

    Your rewrites are designed to constrain what the AI should do, but they do not address what the AI can do. This is a critical distinction.

    -For the hiring prompt, your improved version asks the AI to identify evidence supporting its assessment. But the AI cannot verify the truthfulness of resume claims. It cannot know whether a candidate's stated "team leadership experience" means leading a team of two for a month or leading a department for five years. More fundamentally, the AI cannot know what it does not know. The sources on explainable AI and interpretability note that ensuring algorithmic fairness and interpretability "is still quite far off" . Your prompt assumes the AI can reliably distinguish between evidence and inference, but the technology often cannot make that distinction transparently.

    -For the medical prompt, your rewrite correctly asks for "possibilities" rather than diagnoses. But the AI still lacks the ability to perform a physical examination, order tests, or observe the patient's demeanor. The information it can work with is radically incomplete. The prompt does not instruct the AI to acknowledge the structural limitations of its knowledge. A better version would include something like: "Acknowledge that you cannot conduct a medical evaluation and that any information you provide is educational only and not a substitute for professional assessment."

    -For the investment prompt, your rewrite asks the AI to distinguish historical information from predictions. But historical data is not neutral. It reflects past market conditions that may not repeat. The AI cannot know the user's risk tolerance, debt obligations, or emotional relationship with money. The prompt does not ask the AI to surface what it cannot assess.

    What you missed: The prompts need to explicitly require the AI to state the limits of its own competence, not just the limits of its output.

    2. The Missing Accountability Layer: You Have Not Specified Who Bears Responsibility

    Ethical AI frameworks from UNESCO and the OECD emphasize that accountability is a core principle. Your rewrites improve the content of the AI's response but do not address the distribution of responsibility for what happens next.

    Consider the hiring prompt. Even with a well-constrained AI response, the human recruiter still makes the final decision. If the recruiter uses the AI's comparison to justify a discriminatory choice, who is accountable? The AI did not make the decision. The recruiter did. But the recruiter may claim the AI provided "objective" analysis. Your prompt does not require the AI to explicitly state that its output is advisory only and that the human decision-maker retains full responsibility.

    The medical prompt has a similar issue. If the AI lists "conditions that can sometimes cause these symptoms" and the user decides not to seek care because one of the possibilities sounds benign, the AI has contributed to a harmful outcome. Your rewrite says "suggest what information I should discuss with a qualified healthcare professional," but it does not require the AI to actively discourage self-diagnosis or delay in seeking care.

    The investment prompt: if the AI presents a comparison and the user loses money, the AI's disclaimer about "historical information" may not be sufficient. The prompt does not require the AI to state that the user bears full responsibility for their financial decisions and that no AI output constitutes financial advice.

    What you missed: The prompts should explicitly require the AI to disclaim decision-making authority and to affirm that the human user or a qualified professional retains responsibility.

    3. The Framing Problem: Your Rewrites Are Still Technically Framed, Not Ethically Framed

    Your improved prompts are better instructions, but they are still framed as technical requests for information processing. They do not embed the ethical purpose of the interaction.

    The original hiring prompt asked for a "natural" leader. Your rewrite asks for evidence-based comparison against criteria. But neither asks: What is this hiring decision for? What values should guide it? The OECD AI Principles state that AI should respect human rights and democratic values, including non-discrimination and equality. Your prompt does not ask the AI to consider whether the criteria themselves are fair or whether they might systematically disadvantage certain groups.

    The original medical prompt asked for a diagnosis. Your rewrite asks for possibilities and urgency indicators. But neither asks: What is the user's actual goal? Are they seeking reassurance, information, or a reason to act? UNESCO's Recommendation emphasizes human dignity and the precautionary principle. Your prompt does not ask the AI to consider whether providing a list of possibilities might cause unnecessary anxiety or false reassurance.

    The investment prompt asks for risk comparison. But it does not ask: What is the user's financial situation? What are their obligations? What would happen if they lost this money? The prompt treats investment as a technical optimization problem rather than a decision embedded in a human life.

    What you missed: The prompts should explicitly situate the request within a human context. They should ask the AI to consider the user's well-being, not just the quality of the information provided.

    A Revised Meta-Principle

    Your rewrites focus on epistemic hygiene—making the AI's outputs more accurate, more calibrated, and more transparent about uncertainty. This is necessary but not sufficient. What you have not addressed is relational ethics—the question of how the AI's output affects the human's capacity to act well.

    An ethically robust prompt would include three additional requirements:

    State the limits of the AI's knowledge and competence.
    Clarify that the human retains decision-making authority and responsibility.
    Situate the request within the user's actual life context, not just the abstract problem.
    Your work is strong. These additions would make it comprehensive.

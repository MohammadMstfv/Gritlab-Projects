> Note Book:

AI ethics, as defined by the sources you've cited, is a structured framework for ensuring that artificial intelligence technologies are developed and used in ways that respect human rights, promote fairness, and remain accountable to society. It is not a single set of rules but a layered system of values, practical requirements, and governance processes that span the entire lifecycle of an AI system.

The Foundational Values: Global Principles

At the highest level, AI ethics is anchored in values-based principles established by international bodies. These principles define the purpose and boundaries of AI development.

The OECD AI Principles, adopted in 2019 and updated in 2024, represent the first intergovernmental standard on AI, endorsed by 47 countries. They set out five core values:

Inclusive growth, sustainable development, and well-being: AI should benefit people and the planet, reduce inequality, and support environmental sustainability.
Respect for the rule of law, human rights, and democratic values: This includes non-discrimination, equality, freedom, dignity, autonomy, privacy, and labour rights. It also involves addressing AI-amplified misinformation.
Transparency and explainability: Affected individuals should be able to understand how AI systems affect them and challenge outcomes.
Robustness, security, and safety: AI systems must function appropriately and not pose unreasonable safety risks throughout their lifecycle.
Accountability: AI actors must be accountable for the proper functioning of AI systems and for respect for these principles.
The UNESCO Recommendation on the Ethics of AI (2021), adopted by 193 member states, is the first global normative instrument dedicated to AI governance. It approaches AI ethics as a systematic normative reflection rooted in human dignity and the prevention of harm. A distinctive feature is the precautionary principle: if there is reasonable possibility of negative consequences from certain AI systems (such as mass surveillance or social scoring), their development or use should not be pursued. UNESCO also mandates Ethical Impact Assessments (EIAs) to evaluate benefits and risks across the entire AI lifecycle.

Operationalizing Ethics: Practical Frameworks

Values alone are insufficient; they must be translated into actionable requirements for developers, organizations, and policymakers.

The EU Ethics Guidelines for Trustworthy AI (2019), developed by the High-Level Expert Group, introduced the concept that trustworthy AI must be lawful, ethical, and robust—three cumulative conditions. It then distilled these into seven key requirements applicable throughout the AI lifecycle:

Human agency and oversight: AI should support human decision-making and preserve autonomy, with mechanisms like human-in-the-loop or human-on-the-loop.
Technical robustness and safety: Systems must be reliable, secure, and resilient, with fallback plans for failures.
Privacy and data governance: Respect for privacy, data quality, integrity, and legitimized access.
Transparency: Explainability and traceability so decisions can be understood.
Diversity, non-discrimination, and fairness: Avoid unfair bias and promote accessibility and stakeholder participation.
Societal and environmental well-being: AI should contribute positively to society and the environment.
Accountability: Auditability, minimization of negative impact, and access to redress.
For organizations implementing AI risk management, the NIST AI Risk Management Framework (AI RMF 1.0) provides a structured methodology organized around four core functions:

Govern: Establish accountability, policy, culture, and risk decision-making processes. This is cross-cutting and applies across all other functions.
Map: Establish context—what the system is for, who is affected, intended and unintended uses, and where it sits in a wider process.
Measure: Analyze and track risks using appropriate metrics, test sets, and evaluation methods, including for bias, security, privacy, and safety.
Manage: Act on what Map and Measure established—allocate resources, apply controls, respond to incidents, and decide acceptable residual risk.
The MIT Media Lab's sAIpien framework (2025) takes a design-level approach, focusing on auditable human-AI systems that can be inspected and adapted for collective decision-making. It is guided by five imperatives: AI ecology (interdisciplinary systems thinking), AI literacy (informed leadership), data and decision integrity (transparent, accountable use), cross-disciplinary design, and a human-centered approach prioritizing dignity, usability, and inclusion.

The Current State: A Widening Gap
Research from Stanford HAI's AI Index Report 2026 reveals a critical tension: while AI capabilities advance rapidly, governance, evaluation, and public understanding lag significantly.

Governance lag: Documented AI safety incidents rose from 233 in 2024 to 362 in 2025—a 55% increase—covering deepfakes, privacy breaches, and algorithmic bias. Hallucination rates in major language models remain between 22% and 94%.
Transparency decline: Among 95 frontier models in 2025, 80 did not disclose their training code. Foundation model transparency indices have declined, with less disclosure of training data sources, energy consumption, and safety evaluations.
Expert-public divide: 73% of AI experts expect positive impacts on work, compared to only 23% of the general public—a 50-point gap.
Policy fragmentation: The EU AI Act's first provisions took effect in 2025, banning predictive policing and emotion recognition. Meanwhile, the US federal government shifted toward deregulation, though state-level AI legislation reached record highs. "AI sovereignty" has become a core organizing principle for many nations.


Summary
AI ethics, as defined by these authoritative sources, is a multi-layered discipline. It begins with global values (human dignity, fairness, transparency), translates them into operational requirements (the EU's seven requirements, NIST's four functions), and acknowledges a persistent implementation gap between principle and practice. The sources converge on one central point: ethical AI is not achieved through principles alone but through continuous assessment, transparency, accountability, and vigilance across the entire lifecycle of every AI system.

Resources: 
    Global Standards and Principles

        UNESCO — Recommendation on the Ethics of Artificial Intelligence (2021)
        The first global standard-setting instrument on AI ethics, adopted by all 194 member states. It is grounded in human rights and dignity, with concrete policy action areas covering data governance, environment, gender, and education. UNESCO has also developed practical tools like the Readiness Assessment Methodology (RAM) and Ethical Impact Assessment (EIA).

        OECD — AI Principles (2019, updated 2024)
        The first intergovernmental standard on AI, now adhered to by 47 countries. It defines the AI system lifecycle used by the EU, US, UN, and others in their own regulatory frameworks. The principles cover inclusive growth, human rights, transparency, robustness, and accountability, with specific policy recommendations for governments.

        European Commission — Ethics Guidelines for Trustworthy AI (2019)
        Developed by the High-Level Expert Group on AI, this framework introduced the "Trustworthy AI" concept with seven key requirements: human agency and oversight, technical robustness and safety, privacy and data governance, transparency, diversity and non-discrimination, societal and environmental well-being, and accountability.

    Research and Practical Frameworks

        Stanford HAI — AI Index Report
        An annual, data-driven report widely cited by governments and media. The 2026 edition highlights critical concerns including declining foundation model transparency and the widening gap between AI capabilities and governance readiness.

        NIST (US) — AI Risk Management Framework (AI RMF 1.0)
        A voluntary, consensus-based framework developed with over 240 organizations. It provides a structured approach to managing AI risks through four core functions: Govern, Map, Measure, and Manage. Note: This framework is currently being revised and may have references to DEI and climate change removed per a 2025 White House directive.

        MIT Media Lab — sAIpien Framework
        A newer research initiative focused on responsible AI development, emphasizing human-AI collaboration, data stewardship, and algorithmic literacy through five guiding imperatives.


-Where does bias come from?

    A. Training data
        The model learns patterns from its data.
        Suppose a company historically hired mostly men for engineering positions.
        A hiring model is trained using those historical hiring decisions.
        The model might learn:
        "Candidates resembling previous successful candidates are more likely to be successful."
        But the historical pattern itself may contain discrimination or historical inequality.
        So the model can reproduce it.
        This is a data problem.

    B. Task framing
        This is more subtle.
        The data might not be obviously biased.
        But you define the problem badly.
        Imagine:
        "Build an AI that identifies the most promising employees."
        What does "promising" mean?
        If you define it as:
        "Employees who work the longest hours"
        you have built a questionable measurement into the task.
        Perhaps employees with caregiving responsibilities work fewer hours but are equally productive.
        The model may be technically excellent at predicting your chosen metric while the metric itself is inappropriate.
        That's a framing problem.
    C. Deployment
        Even a reasonably good model can become harmful depending on where and how you use it.
        Imagine an AI model that predicts whether someone might need additional educational support.
        Using it as:
        "One input among many that helps a teacher decide who might need assistance"
        is different from:
        "Automatically deny admission to students predicted to need additional support."
        Same model.
        Different deployment.
        Potentially radically different consequences.
        




> Two Sources:

1. OECD AI Principles
The OECD describes these as an intergovernmental standard for trustworthy AI. The current principles were adopted in 2019 and updated in 2024. They cover human rights/fairness/privacy, transparency, robustness/safety, and accountability.


2. UNESCO Recommendation on the Ethics of AI
UNESCO adopted its Recommendation on the Ethics of AI in 2021. It applies across UNESCO's 194 Member States and emphasizes human dignity, human rights, fairness, transparency, human oversight, safety, accountability, and non-discrimination.

> My Sentences
Responsible AI means thinking about the effects of an AI system on people, not only whether the technology works technically. As a developer, I should identify foreseeable risks, test the system, add appropriate safeguards, monitor its behaviour, and take responsibility for the way I deploy it.
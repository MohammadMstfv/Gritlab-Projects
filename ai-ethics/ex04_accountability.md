1. Responsibility is shared but asymmetric.
    As the developer who shipped the feature, I bear primary responsibility, alongside my company. The model provider has secondary responsibility. The user has limited responsibility, and only if they ignored clear warnings or misused the system.

    > Why “the model did it” is not acceptable

    An LLM is not a moral agent. It has no intent, no duties, and no capacity to be accountable. It is a component — like a database, a calculator, or a network library. If I ship a feature that relies on it, I am responsible for understanding its failure modes, testing it, constraining it, and monitoring it. “The model did it” is an abdication of engineering duty.

    The NIST AI Risk Management Framework is built on the premise that organizations must Govern, Map, Measure, and Manage AI risks. The EU Ethics Guidelines for Trustworthy AI require human agency, oversight, and accountability. The OECD AI Principles state that AI actors must be accountable for the proper functioning of AI systems. None of these frameworks allow responsibility to be transferred to the tool itself.

    If I cannot explain how the system could fail, how I tested for that failure, and what safeguards prevent harm, I should not ship it. The model’s behavior is foreseeable enough that engineers are expected to design around it — especially for high-stakes use cases.

    > Who is responsible, and why

    Me: Primary, I chose the model, wrote the prompts, designed the integration, ran the tests (or failed to), and shipped the feature. I had proximate control over whether it was safe enough to release.

    My company:	Primary / joint, It owns the product, profits from it, sets risk tolerance, allocates resources for safety, and bears legal liability. It is responsible for governance, review, and incident response.

    Model provider:	Secondary, It built the base model and knows its limitations. It should provide accurate documentation, safety evaluations, and clear warnings about known failure modes. But it cannot foresee every downstream use.

    User: Limited, Users can be responsible for verifying high-stakes information or not misusing a tool. But they cannot inspect the model, and they may reasonably trust a product that presents itself as reliable. Blaming the user for a hidden defect is usually wrong.

    > Two concrete things I would put in place before shipping

    1. *Domain-specific risk assessment and pre-deployment evaluation.*
    I would define what “harm” means for this feature, identify high-risk failure modes (e.g., confident false statements, missing context, harmful recommendations), and build a test suite that measures hallucination rates, calibration, and edge cases. I would use red-teaming with domain experts and affected users. For high-stakes outputs, I would set explicit thresholds: if the model cannot meet them, the feature does not ship, or it ships only with mandatory human review. This is the Measure and Manage functions of the NIST AI RMF.

    2. *Runtime guardrails and human-in-the-loop escalation.*
    I would not rely on a disclaimer. I would constrain the model’s scope, ground it in vetted sources where possible (e.g., retrieval-augmented generation with citations), and require it to refuse or escalate when evidence is insufficient or confidence is low. For any consequential action — medical, financial, legal, employment — a human must review before the action takes effect. I would also add audit logging, monitoring for drift and incidents, a rollback plan, and a kill switch. This aligns with the EU Trustworthy AI requirements for human oversight, transparency, and accountability.

    >Going further: one argument for and against each party’s responsibility

    Developer
    For: You designed and shipped the feature; you had the most direct ability to prevent the harm.
    Against: You may have been pressured by deadlines, lacked access to model internals, or relied on provider claims that turned out to be false.
    
    Company
    For: It controls resources, sets priorities, and is legally liable for the product it sells.
    Against: If leadership ignored engineering warnings or overrode safety concerns, blame may shift upward — but the company as an entity still bears responsibility.
    
    Model provider
    For: It built the model, knows its weaknesses, and has a duty to document and mitigate known risks.
    Against: It cannot control how downstream developers deploy the model or what prompts they use.
    
    User
    For: Users should exercise judgment and not blindly trust AI for high-stakes decisions.
    Against: Users cannot audit the model, and a polished product may reasonably create trust. Harm from a hidden defect is not the user’s fault.
  

    *I am responsible because I shipped it. My company is responsible because it put it into the world. The model provider is responsible for what it knew and disclosed. The user is responsible only at the margins. The correct engineering mindset is: assume the model will fail, and design so that its failure does not cause harm. “The model did it” is not an answer — it is a confession that I did not do my job.*
A:
    1. Found Concrete Example

        > Amazon's experimental recruiting system.
        

    2. Explain in My Words: where the bias most likely originates: the training data, the way the task was framed, or how the system is deployed. Justify your choice.

        Since it worked on historical data, and most of the employments where male in earlier times, the bias is originated from the data. But also the main problem originates from the historical decisions. The AI is using it's only data available but the data is in this way, because it was this way in earlier times.

    3. Who is harmed

        First of all, women who look for employment and face this AI before any human interview, because they will not be able to surpass the first lvl which would be AI review. But also the other candidates whose characteristics correlate with historical patterns, employers because they may lose qualified candidates and the organization itself through discrimination/legal/reputational consequences


B:

AI Stress Test

    You identified the training data correctly, but you missed a distinct cause: the way the system was deployed and the nature of the feedback loop it created.

    The system was not merely a passive reflection of historical data. It was actively deployed as a screening tool that recruiters used to find candidates . The bias was not just in the data; it was operationalized through a mechanism that gave candidates scores from one to five stars . This deployment created a self-reinforcing feedback loop: the biased system would disproportionately reject women, which meant future hiring data would continue to show men succeeding, which would further "confirm" the model's bias. The sources note that even after Amazon edited the program to neutralize specific terms like "women's," the engineers admitted there was "no guarantee" the machine wouldn't devise other ways to sort candidates discriminatorily . This shows the bias was not a static data problem but an emergent property of the deployed system.

    Who You Overlooked: The Engineering Team and the "Holy Grail" Mindset

    You listed women, other candidates, employers, and the organization. You missed the engineering team itself and the culture that produced the system.

    The team spent years (2014–2017) building 500 computer models, only to see the project fail and the team disband . The sources describe a culture of hype—"Everyone wanted this holy grail"—where the goal was to "spit out the top five" resumes automatically . The engineers were harmed by the professional failure and the realization that their technical work had embedded discrimination. The sources also quote experts like Nihar Shah from Carnegie Mellon, noting that ensuring algorithmic fairness and interpretability "is still quite far off" . This points to a systemic harm: the overpromising of AI capabilities without the technical safeguards (like fairness auditing or interpretability) needed to prevent harm. The "holy grail" mindset itself is an overlooked cause, because it prioritizes automation speed over ethical checks and human oversight.


C:
Reflection

    I should think myself and oversee such things. Although, us humans have our own biases and should be careful about them, in anything we do and develop.
    I try to be a better person myself first. And then yes, I must think of these flaws beforehand. I am responsible for deciding what the facts mean in the context of the system I'm building.
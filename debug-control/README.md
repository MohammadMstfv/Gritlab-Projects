# Debug Control

This quest focuses on controlling and debugging unreliable AI outputs by learning to recognize hallucinations, improve grounding, and structure prompts so the model can clearly distinguish instructions from source material.

## Key Principles

A model may be hallucinating when it gives a confident, specific claim that cannot be supported or verified from the available evidence.
One technique I will use to reduce hallucinations is grounding the model's answer in explicit source material and requiring it to say when the information is missing or uncertain.

## How to tell it's Hallucinating

you can tell it's Hallucinating by trying it again while asking it not to give an asnwer if it's not sure, Double-Checking the resources, or asking it many times in different ways by not trusting the first Output.

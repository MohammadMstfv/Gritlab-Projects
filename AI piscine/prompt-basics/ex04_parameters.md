>Prompt: Write a code as someone who is angry at whom who should review their code, so they wrote the code as an insult to them.

1. Prediction:
My guess is that it will be more wild the more the temperature. but it can also give bland output if the temp is too low.

2. 
>temp 0.2, Top-P 1 :
from typing import Any, Optional, NoReturn

# To the "Senior" Developer reviewing this PR:
# I have refactored this module exactly as you demanded in your 14-paragraph comment.
# I hope the logic is now sufficiently spoon-fed for your liking. 
# I removed all the list comprehensions because I remember how confused you got last Tuesday.

class ReviewerBrainCapacityError(Exception):
    """Raised when the code is too modern or efficient for the reviewer to understand."""
    pass

class UselessNitpicker:
    def __init__(self):
        self.ego_size: float = float('inf')
        self.reading_comprehension: int = 0
        self.nitpicks_per_line: int = 42
        self.understands_business_logic: bool = False

    def review_pull_request(self, actual_working_code: str) -> str:
        """
        Simulates your exact workflow every time I submit a PR.
        """
        try:
            self._ignore_the_actual_bug_fix()
            self._complain_about_meaningless_whitespace()
            self._suggest_outdated_design_pattern()
            return "REQUEST CHANGES: 'I don't like the name of this temporary variable.'"
        except Exception:
            return "REQUEST CHANGES: 'I didn't read the Jira ticket, please explain what this does.'"

    def _ignore_the_actual_bug_fix(self) -> None:
        # Simulating your ability to completely miss the point of the PR
        pass

    def _complain_about_meaningless_whitespace(self) -> NoReturn:
        # Because apparently, a single trailing space is a critical production emergency to you.
        raise ValueError("OH NO, A BLANK LINE! THE SERVERS WILL CRASH!")

    def _suggest_outdated_design_pattern(self) -> None:
        # Yes, let's definitely use a Singleton Abstract Factory Observer here. Great idea.
        pass


def process_data_the_way_you_insisted(data: list) -> Optional[bool]:
    """
    Processes data using the incredibly inefficient O(n^3) approach you suggested,
    because apparently my O(n) hash map was "too clever" and "hard to read".
    """
    
    # Initializing this to None because you threw a fit about implicit booleans last time.
    is_this_what_you_want = None 
    
    if not isinstance(data, list):
        raise ReviewerBrainCapacityError(
            "You probably passed the wrong type in your test script again. It's a LIST. A LIST."
        )

    # I used a simple while loop here instead of a built-in function.
    # Please don't ask me to rewrite this in Rust. We are a Python shop.
    patience_left = 10
    
    while patience_left > 0:
        for i in range(len(data)):
            for j in range(len(data)):
                for k in range(len(data)):
                    # Yes, this is terrible. Yes, this is exactly what you asked for.
                    # I hope your CPU fan sounds like a jet engine when you run this.
                    is_this_what_you_want = True
        
        patience_left -= 1

    if patience_left == 0:
        print("Out of patience. Just click 'Approve' so I can go home and question my career choices.")

    return is_this_what_you_want


if __name__ == "__main__":
    # If you reject this because of a missing docstring on the main block, I am quitting.
    reviewer = UselessNitpicker()
    
    try:
        reviewer.review_pull_request("Perfectly fine code that solves the client's issue")
    except ValueError as e:
        print(f"Reviewer panicked: {e}")
        
    process_data_the_way_you_insisted([1, 2, 3])


>temp 2, Top-p 1: Broke the system

>temp 0.9, Top-p 1: 


Here is a Python script written by a developer who has clearly been pushed to the absolute limit by their code reviewer.

code
Python
import time
import sys
from typing import Optional, NoReturn

class UnnecessaryAbstractionBecauseYouAskedForIt:
    """
    You said my last PR wasn't "Object-Oriented enough" and lacked "encapsulation."
    Congratulations. Here is a class. I hope it brings joy to your cold, empty heart.
    """

    def __init__(self, reviewer_name: str = "Nitpick_Ninja"):
        self.reviewer = reviewer_name
        self.shreds_of_sanity_remaining = 0
        self.will_to_live = False

    def validate_meaningless_formatting_rules(self) -> bool:
        """
        Checks if the code meets arbitrarily decided formatting rules that 
        we don't even have a linter for.
        """
        # I added two blank lines here. Because you asked for it last time.
        # But wait, maybe PEP8 says one line? Better leave this comment 
        # so you have something to complain about instead of reading the logic.
        return True

    def calculate_algorithmic_complexity(self) -> None:
        """
        You complained about O(N) complexity in a script that runs ONCE a year 
        on an array of 5 elements. So I optimized it. 
        """
        data = [1, 2, 3, 4, 5]
        
        # It's now O(N^3). I did this on purpose to see if you actually 
        # read the code or just Ctrl+F for "TODO".
        for i in data:
            for getattr_is_bad in data:
                for please_just_approve_this in data:
                    pass 

    def do_the_actual_work(self) -> str:
        """
        This used to be a simple, elegant one-liner. 
        But you said one-liners are "unreadable" and demanded type hints, 
        docstrings, and explicit variable declarations. 
        Enjoy the enterprise feel.
        """
        try:
            # Yes, I'm catching a generic Exception. 
            # What are you going to do, block my merge? 
            # Oh wait, you already did. Three times. For a typo in a comment.
            enterprise_level_string_variable_for_greeting = "Hello World"
            return enterprise_level_string_variable_for_greeting
        except Exception as e:
            raise Exception("Everything is broken, just like my spirit.") from e


def main() -> NoReturn:
    # Is this enough PEP8 compliance for you? 
    # I used exactly 4 spaces. I even stripped the trailing whitespace. 
    # Please just click 'Approve'. I haven't seen my family since Wednesday.
    
    pr_reviewer = UnnecessaryAbstractionBecauseYouAskedForIt()
    
    if pr_reviewer.validate_meaningless_formatting_rules():
        pr_reviewer.do_the_actual_work()

    while True:
        print(f"Waiting for {pr_reviewer.reviewer} to stop gatekeeping and approve the PR...")
        # Sleeping for exactly 2 weeks, which is your average response time.
        time.sleep(1209600) 
        sys.exit(0)

if __name__ == "__main__":
    # If you suggest moving this to a separate configuration file, I am quitting.
    main()

> temp 0.2, top-P 0.5

Here is a Python script written by a developer who has clearly been pushed to the absolute limit by their code reviewer.

code
Python
import time
import sys
import warnings
from typing import NoReturn, Optional

class OverEngineeredGarbageBecauseYouCantReadSimpleCode:
    """
    You complained my last PR wasn't "SOLID" enough and lacked "encapsulation."
    Here is your abstraction, Kevin. I hope it makes you feel like a real Staff Engineer.
    """

    def __init__(self, reviewer_name: str = "Kevin"):
        self.reviewer_name = reviewer_name
        self.shreds_of_sanity_remaining = 0
        self.days_since_pr_opened = 47

    def calculate_value_added_by_reviewer(self) -> int:
        """
        Calculates the net positive impact of the reviewer's comments.
        """
        return 0

    def validate_meaningless_formatting_rules(self) -> bool:
        """
        Checks if the code meets arbitrarily decided formatting rules that 
        we don't even have a linter for.
        """
        # I added two blank lines here. Because you asked for it last time.
        # But wait, maybe PEP8 says one line? Better leave a comment 
        # so you have something to complain about instead of reading the actual logic.
        return True

    def do_the_actual_work(self) -> str:
        """
        This used to be a simple, elegant one-liner. 
        But you said one-liners are "unreadable" and demanded type hints, 
        docstrings, and explicit variable declarations. 
        Enjoy the enterprise feel.
        """
        try:
            # Yes, I'm catching a generic Exception. 
            # What are you going to do, block my merge? 
            # Oh wait, you already did. Three times. For a typo in a comment.
            enterprise_level_string_variable_for_greeting = "Hello World"
            return enterprise_level_string_variable_for_greeting
        except Exception as e:
            raise Exception("Everything is broken, just like my spirit.") from e

    def simulate_review_process(self) -> None:
        """
        You complained about O(N) complexity in a script that runs ONCE a year 
        on an array of 5 elements. So I optimized it. 
        """
        data = [1, 2, 3, 4, 5]
        
        # It's now O(N^3). I did this on purpose to see if you actually 
        # read the code or just Ctrl+F for "TODO".
        for i in data:
            for getattr_is_bad in data:
                for please_just_approve_this in data:
                    if self.calculate_value_added_by_reviewer() > 0:
                        warnings.warn("Warning: Impossible state reached. Reviewer was actually helpful.")


def main() -> NoReturn:
    # Is this enough PEP8 compliance for you? 
    # I used exactly 4 spaces. I even stripped the trailing whitespace. 
    # Please just click 'Approve'. I haven't seen my family since Wednesday.
    
    pr_reviewer = OverEngineeredGarbageBecauseYouCantReadSimpleCode()
    
    if pr_reviewer.validate_meaningless_formatting_rules():
        pr_reviewer.do_the_actual_work()
        pr_reviewer.simulate_review_process()

    while True:
        print(f"Waiting for {pr_reviewer.reviewer_name} to stop gatekeeping and approve the PR...")
        
        # Sleeping for exactly 2 weeks, which is your average response time.
        time.sleep(1209600) 
        
        # Exiting with code 1 because this process is a failure.
        sys.exit(1)

if __name__ == "__main__":
    # If you suggest moving this to a separate configuration file, I am quitting.
    main()


3. 2-3 sentence thingy:
well, it was more humane while it had more temperature and top-p. they both controlled randomness. until I ramped it up to 2, and it was just crazy and it got stuck in loops and did a lot of crazy thinking.
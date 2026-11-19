1. Ask for answer Only:

    Three people (A, B, C) need to cross a river. The boat holds at most two. A and B cannot be left alone together, B and C cannot be left alone together. How do they all cross? Give only the sequence of crossings.

    > A and B cross, A returns, A and C cross, A returns, A and B cross

2. Ask with STEP BY STEP:

    Three people (A, B, C) need to cross a river. The boat holds at most two. A and B cannot be left alone together, B and C cannot be left alone together. How do they all cross?
    List every possible legal move and eliminate invalid ones until you reach a complete solution.
    
    >Under the usual river-crossing rules, there is no legal way for A, B, and C to all cross. The puzzle as stated is impossible.

        Let “valid” mean: on either bank, A and B may not be together without C, and B and C may not be together without A.

        Start:

        Left bank: A, B, C
        Right bank: empty
        Boat: left
        All possible first moves from left to right

        A alone → Left: B, C.
        Invalid: B and C are left alone together.
        B alone → Left: A, C; Right: B.
        Valid.
        C alone → Left: A, B.
        Invalid: A and B are left alone together.
        A + B → Right: A, B.
        Invalid: A and B are on the right alone together.
        A + C → Left: B; Right: A, C.
        Valid.
        B + C → Right: B, C.
        Invalid: B and C are on the right alone together.
        So the only legal first moves are:

        B crosses alone, or
        A and C cross together.
        Check the first legal branch: B crosses alone

        B crosses: Left = A, C; Right = B.
        Boat is now on the right. Only B is there, so B must return.
        B returns: back to the starting state.
        This makes no progress, so it cannot lead to a solution.

        Check the second legal branch: A and C cross together

        A and C cross: Left = B; Right = A, C.
        Boat is on the right. Someone must return.
        Possible return moves:

        A returns: Left = A, B; Right = C.
        Invalid: A and B are left alone together.
        C returns: Left = B, C; Right = A.
        Invalid: B and C are left alone together.
        A and C return together: back to the starting state.
        So this branch also cannot continue legally.

        Conclusion

        All legal first moves have been exhausted. Every possible branch either immediately becomes invalid or returns to the starting state.

        Therefore, A, B, and C cannot all cross the river under the stated rules. One of the constraints must be changed or removed for a solution to exist.


3. did the step-by-step version change the answer or only the presentation? On what kind of problem would you expect CoT to actually change correctness?

    Well in most cases it gave the correct answer anyways. But after trying many things on many agents, I did get a wrong answer once, but when it did it step-by-step it got it correct. 
    I expect it to change the asnwer for trick questions or such. Because when it wants to give you steps, it gives you the steps and that is its answer, so it can say there's no answer at all or something. but when you want direct answer, it should give you an answer so it gives something wrong.
1. Weak prompt
    
    summarize this

    func CountIf(f func(string) bool, tab []string) int {
        count := 0
    for _, v := range tab {
    if f(v) {
    count++
            }
        }
    return count
    }

    > Answer:

    This Go function CountIf counts how many strings in a slice (tab) satisfy a given condition.

    It takes a predicate function f (which returns true/false for a string) and a string slice.
    It loops through the slice, applies f to each element, and increments a counter whenever f returns true.
    It returns the final count.

2. Improved Prompt

    Summarize the following in 3 bullet points: what is it about, what is it used for, what is a problem about it
    
    func CountIf(f func(string) bool, tab []string) int {
        count := 0
    for _, v := range tab {
    if f(v) {
    count++
            }
        }
    return count
    }"

    > Answer

        -What it is about: A Go function named CountIf that takes a predicate function f (of type func(string) bool) and a string slice tab, then returns the number of elements in the slice for which f returns true.

        -What it is used for: Counting how many strings in a slice satisfy a custom condition (e.g., length checks, pattern matching, or any boolean test), acting as a simple higher-order “count if” utility.

        -What is a problem about it: The provided code has poor/inconsistent indentation (the for loop and if body are misaligned), which reduces readability and can cause maintenance issues; it also lacks any nil-checks for the function or slice.

3. Run on Two Input

    1. the thing you saw and the answer is up here.

        -What it is about: A Go function named CountIf that takes a predicate function f (of type func(string) bool) and a string slice tab, then returns the number of elements in the slice for which f returns true.

        -What it is used for: Counting how many strings in a slice satisfy a custom condition (e.g., length checks, pattern matching, or any boolean test), acting as a simple higher-order “count if” utility.

        -What is a problem about it: The provided code has poor/inconsistent indentation (the for loop and if body are misaligned), which reduces readability and can cause maintenance issues; it also lacks any nil-checks for the function or slice.

    2. I gave it a document and here is the answer

        -What it is about: A university final-project report (in Persian) describing the design of a dynamic/adaptive music system for a computer game. It covers differences between AAA and indie game music, a case study of Hades, a simple original game story/structure with four distinct spaces, specific dynamic-music features (layer cut-offs during the final boss fight, music changes based on environment or player health, two-phase boss music), and an AI-generated audio-visual prototype.

        -What it is used for: Academic documentation of research, design decisions, and an early prototype (four music tracks + concept art + showcase video) to support later full implementation in Unity, serving as the student’s final bachelor’s project report.

        -What is a problem about it: The main challenges noted are keeping the story simple yet narratively coherent enough to justify musical variety, iteratively refining AI-generated (Suno) music to match the intended identity of each space, and deciding between Unity’s built-in audio tools versus specialized middleware (FMOD/Wwise) for the technical implementation still to come.

4. One Sentence: did the improved prompt give a consistent, useful shape both times?

    Yes, it gave the same shape of answer. eventhough when it needed more.


The buggy version sorts in descending order simply because of a small mistake. 
I changed the direction of a bigger/smaller


ChatGPT:
If you're referring to the buggy Python function above, changing < back to > fixes the sorting direction. After that, the function is basically correct, but there are still some edge cases and design considerations.
(Here, it gave me some examples, I ran them and there was no problem with the code)

Large positive/negative integers should also work because Python integers don't have Go's fixed integer overflow issue.

One important design point
Your function doesn't rearrange the nodes themselves. It only swaps their Data:
i.Data, j.Data = j.Data, i.Data
That's perfectly fine if the exercise allows it.
For example:
Node A: Data=5 → Node B: Data=2
becomes:
Node A: Data=2 → Node B: Data=5
The nodes stay in exactly the same positions.
Is there a cleaner approach?
Yes. Python already has a very clean sorting mechanism, but it doesn't directly apply to your custom linked list. You could first convert the list to a Python list, sort it, and put the values back:

def ListSort(l):
    if l is None:
        return None

    values = []

    current = l
    while current is not None:
        values.append(current.Data)
        current = current.Next

    values.sort()

    current = l
    for value in values:
        current.Data = value
        current = current.Next

    return l

However, for your Go/Piscine exercise, I'd keep your original approach. The point is probably to practice linked-list traversal and sorting rather than rely on Python's built-in .sort().
The main thing I'd want you to test now is duplicates, negative numbers, empty list, one node, already sorted, and reverse sorted. 
(Which I did)
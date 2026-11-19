class NodeI:
    def __init__(self, data):
        self.Data = data
        self.Next = None


def ListSort(l):
    if l is None:
        return None

    i = l
    while i is not None:
        j = i.Next
        while j is not None:
            if i.Data < j.Data:  
                i.Data, j.Data = j.Data, i.Data
            j = j.Next
        i = i.Next

    return l
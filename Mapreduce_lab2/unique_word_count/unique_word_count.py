from mrjob.job import MRJob
import re

# Regular expression to extract words from input file
WORD_RE = re.compile(r"\b\w+\b")

class Word_Count_MRJob(MRJob):

    def mapper(self, _, line):
        words = WORD_RE.findall(line.lower())
        for word in words:
            yield word, 1

    def reducer(self, word, counts):
        yield word, sum(counts)

if __name__ == "__main__":
    Word_Count_MRJob.run()

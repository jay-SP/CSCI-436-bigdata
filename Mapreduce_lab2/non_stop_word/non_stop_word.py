#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Mon Jul 31 15:09:59 2023

@author: jp
"""

from mrjob.job import MRJob
import re

#regex
WORD_RE = re.compile(r"\b\w+\b")

#Stop_words
STOP_WORDS = {"the", "and", "of", "a", "to", "in", "is", "it"}

class NonStopWordCount(MRJob):

    def mapper(self, _, line):
        words = WORD_RE.findall(line.lower())
        for word in words:
            if word not in STOP_WORDS:
                yield word, 1

    def reducer(self, word, counts):
        yield word, sum(counts)

if __name__ == "__main__":
    NonStopWordCount.run()

#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Mon Jul 31 19:26:31 2023

@author: jp
"""

from mrjob.job import MRJob
import re

class Word_Bigram_Count(MRJob):

    def mapper(self, _, line):
        #RE
        WORD_RE = re.compile(r"\b\w+\b")

        words = WORD_RE.findall(line)
        num_words = len(words)

        for i in range(num_words - 1):
            bigram = f"{words[i].lower()},{words[i + 1].lower()}"
            yield bigram, 1

        # Handling the last word as a single word bigram
        if num_words >= 1:
            last_word_bigram = f"{words[num_words - 1].lower()}"
            yield last_word_bigram, 1

    def reducer(self, bigram, counts):
        yield bigram, sum(counts)

if __name__ == "__main__":
    Word_Bigram_Count.run()

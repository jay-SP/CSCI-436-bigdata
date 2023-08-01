#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Mon Jul 31 19:07:48 2023

@author: jp
"""

from mrjob.job import MRJob
import re

# Regular expression to extract words from input lines
WORD_RE = re.compile(r"\b\w+\b")

class Inverted_Index(MRJob):

    def mapper(self, _, line):
        # Extract the document ID and document content from the input line
        parts = line.split(":", 1)
        if len(parts) != 2:
                #(invalid input)
            return

        document_id, document = parts
        document_id = document_id.strip()

        # Tokenizing the words in the document
        words = WORD_RE.findall(document.lower())

        # key-value pairs with word as key and document ID as value
        for word in words:
            yield word, document_id

    def reducer(self, word, doc_ids):
        # Creating a set to store unique document IDs for each word
        doc_id_set = set(doc_ids)

        # Converting the set to a comma-separated string of document IDs
        doc_id_str = ", ".join(sorted(doc_id_set))

        # Yielding the word and its corresponding document IDs
        yield word, doc_id_str

if __name__ == "__main__":
    Inverted_Index.run()



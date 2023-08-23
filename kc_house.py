#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Aug 22 21:50:20 2023

@author: jp
"""

from pyspark.sql import SparkSession
from pyspark.ml.feature import VectorAssembler
from pyspark.ml.regression import RandomForestRegressor
from pyspark.ml.evaluation import RegressionEvaluator

# Initialize a Spark session
spark = SparkSession.builder.appName("RandomForestRegression").getOrCreate()

# Load your CSV dataset as a Spark DataFrame
# Assuming you have a CSV file named 'data.csv'
data = spark.read.csv('/Users/jp/Downloads/kc_house_data.csv', header=True, inferSchema=True)

# Select features and target variable
features = ['bedrooms', 'bathrooms', 'sqft_living', 'sqft_lot', 'floors', 'waterfront',
            'view', 'condition', 'grade', 'sqft_above', 'sqft_basement', 'yr_built',
            'yr_renovated', 'zipcode', 'lat', 'long', 'sqft_living15', 'sqft_lot15']
target = 'price'

# Assemble features into a vector column
assembler = VectorAssembler(inputCols=features, outputCol='features')
data = assembler.transform(data)

# Split data into training and testing sets
train_data, test_data = data.randomSplit([0.8, 0.2], seed=42)

# Create and train the Random Forest Regressor model
rf_model = RandomForestRegressor(featuresCol='features', labelCol=target, numTrees=100, seed=42)
rf_fit = rf_model.fit(train_data)

# Make predictions on the test set
predictions = rf_fit.transform(test_data)

# Evaluate the model
evaluator = RegressionEvaluator(labelCol=target, metricName='rmse')
rmse = evaluator.evaluate(predictions)
print(f"Root Mean Squared Error: {rmse}")

# Stop the Spark session
spark.stop()

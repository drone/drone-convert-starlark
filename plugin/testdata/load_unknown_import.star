# This extension does not exist.
load("@test:unknown_extension.star", "example_step")

# Absolute load.
load("@test//subpackage:pipelines_extension.star", "example_pipeline")

def main(ctx):
    return example_pipeline("sample", steps = example_step())

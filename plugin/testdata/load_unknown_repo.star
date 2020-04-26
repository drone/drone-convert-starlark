# Relative load.
load("@test//:steps_extension.star", "example_step")

# Invalid repo load.
load("@test2//subpackage:pipelines_extension.star", "example_pipeline")

def main(ctx):
    return example_pipeline("sample", steps = example_step())

def example_pipeline(name, steps):
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": name,
        "steps": steps,
    }

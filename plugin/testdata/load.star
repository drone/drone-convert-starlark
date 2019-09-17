load('single.star', 'single')

def main(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'default'
  }

def main(args):
    with open('/tmp/show-running', 'w') as f:
        f.write('false')
    f.close()


if __name__ == "__main__":
    main("arg")

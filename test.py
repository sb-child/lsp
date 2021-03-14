import getLinks


def main():
    lg = getLinks.LinksGet()
    link_urls = lg.run()
    for link_url in link_urls:
        print(link_url)


if __name__ == '__main__':
    main()

from selenium import webdriver
from bs4 import BeautifulSoup
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import atexit
from flask_caching import Cache

from flask import Flask, request,jsonify

HTML_PARSER = 'html.parser'
ANTI_BOT_PAGE_TEXT = "Sorry, we just need to make sure you're not a robot. For best results, please make sure your browser is accepting cookies."

app = Flask(__name__)
cache = Cache(app, config={'CACHE_TYPE': 'simple'})


class Product():
    def __init__(self, title, description) -> None:
        self._title = title
        self._description = description

class WebScrapper():
    def __init__(self) -> None:
        self._driver = webdriver.Firefox()
        self.visit_page("https://www.amazon.com/")
    
    def check_if_catpcha(self):
        self._soup = BeautifulSoup(self._driver.page_source,HTML_PARSER)
        p_tag = self._soup.find('p',class_ ='a-last')

        if p_tag:
            print(f'tag found: {p_tag}')
            while True:
                print(p_tag.text)
                self._soup = BeautifulSoup(self._driver.page_source,'html.parser')
                p_tag = self._soup.find('p',class_ ='a-last')
                if p_tag == None:
                    print("////////////////////////////////////////CHAPTA SOLVED///////////////////////////")
                    break
        else:
            print("Catpcha Tag not found")
        
    def find_title(self):
        try:
            title_tag = WebDriverWait(self._driver, 1).until(
                EC.presence_of_element_located((By.ID, 'productTitle'))
            )
            title_text = title_tag.text.strip()
            return title_text
        except Exception as e:
            print(f"Title tag not found: {e}")
        return ""

    def find_description(self):
        try:
            product_description_div = WebDriverWait(self._driver, 1).until(
                EC.presence_of_element_located((By.ID, 'productDescription'))
            )
            product_description_text = product_description_div.text.replace("Product Description", '').replace("Amazon.com", '').strip()
            return product_description_text
        except Exception as e:
            print(f"Product description tag not found: {e}")
        return ""

    def parse_page(self) -> Product: 
        title= self.find_title()
        description= self.find_description()

        return Product(title, description)

    def visit_page(self, url):
        self._driver.get(url)
        self._soup = BeautifulSoup(self._driver.page_source,HTML_PARSER)
        self.check_if_catpcha()
    
    def close(self):
        self._driver.quit()


@cache.memoize(timeout=300)  
def scrape_and_cache(url)  -> dict[str, str]:
    web_scrapper.visit_page(url)
    product = web_scrapper.parse_page()
    return {
        "title": product._title,
        "description": product._description
    }

@app.route('/scrape', methods=['GET'])
def scrape():
    url = request.args.get('url')
    if not url:
        return jsonify({"error": "URL parameter is required"}), 400

    cached_product = scrape_and_cache(url)
    print("returning response from cache")
    return jsonify({
        "title": cached_product['title'],
        "description": cached_product['description'],
    })




if __name__ == "__main__":
    web_scrapper = WebScrapper()
    app.run(debug=True, port = 8082)


@atexit.register
def shutdown():
    web_scrapper.close()


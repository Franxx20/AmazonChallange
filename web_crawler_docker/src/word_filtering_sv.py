from flask import Flask, request, jsonify
import spacy
from collections import Counter
from flask_caching import Cache

app = Flask(__name__)
cache = Cache(app,config={'CACHE_TYPE':'simple'})
nlp = spacy.load("en_core_web_sm")



@cache.memoize(timeout=300)
def filter_text_and_count_words(text):
    text = text.lower()
    doc = nlp(text)
    filtered_words = [token.text.lower() for token in doc if token.pos_ in {"NOUN", "ADJ", "VERB"}]
    word_counters = Counter(filtered_words).most_common(10)
    return word_counters

@app.route('/product', methods=['POST'])
def post_product_description():
    data = request.get_json()  
    
    if data and 'productDescription' in data:
        product_description = data['productDescription']
        word_counts = filter_text_and_count_words(product_description)
        print(word_counts)
        
        response = {
            "wordCounts": [{"word": word, "count": count} for word, count in word_counts],
            "message": "Product description received successfully"
        }
    else:
        response = {
            "wordCounts": [],
            "message": ""

        }
    return jsonify(response)

if __name__ == '__main__':
    app.run(host='127.0.0.1', port=8081, debug=True)
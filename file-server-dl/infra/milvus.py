from pymilvus import MilvusClient
import numpy as np

class MilvusConn():
    def __init__(self):
        self.client = MilvusClient("image_embedding.db")
        self.collection_name = "image_feature"
        if not self.client.has_collection(collection_name=self.collection_name):
            self.client.create_collection(
                collection_name=self.collection_name,
                dimension=768,  # The vectors we will use in this demo has 768 dimensions
                id_type='str',
                max_length=1024
            )
    
    def insert(self, vector: np.ndarray, path: str):
        res = self.client.upsert(collection_name=self.collection_name, data={
            "id": path,
            "vector": vector,
            "path": path
        })

    def search_by_vec(self, vector: np.ndarray, minimal_similarity: float = 0.8):
        res = self.client.search(collection_name=self.collection_name, 
                                 data=[vector], 
                                 search_params={
                                     "params": {
                                         "radius": minimal_similarity,
                                         "range_filter": 1
                                     }
                                 })
        return res[0]

conn = MilvusConn()
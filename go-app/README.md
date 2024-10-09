# Projet API S3 Go MinIO

## Utilisation 

Lancer le server minio depuis le répertoire (minio-tools)
```sh
./minio.exe server ../minio-data
```
Un nouveau dossier minio-data sera alors créer et les buckets seront disponibles dedans

Créer ensuite un alias pour le server MinIO 
```sh 
./mc.exe alias set local http://127.0.0.1:9000 minioadmin minioadmin
```

Il est maintenant possible d'utiliser le CLI minio : 

Fonctions disponibles :


### Lister les buckets 
```sh
./mc.exe ls local
```

### Lister les objets d'un bucket
```sh 
./mc ls local/mybucket

ou 

./mc ls path/to/bucketStorage/TargetBucketName
```

### Créer un bucket 
```sh 
./mc mb local/mybucket

ou 

./mc mb path/to/bucketStorage
```

### Ajouter un fichier 
```sh
./mc cp sample1.txt local/mybucket/

ou 

./mc cp file.txt path/to/bucketStorage/TargetBucketName
```

### Download un fichier 
```sh 
./mc cp local/mybucket/sample1.txt .

ou 

./mc cp path/to/bucketStorage/TargetBucketName/TargetFileName .

```

### Supprimer un fichier 
```sh 
./mc rm local/mybucket/sample1.txt .

ou 

./mc rm path/to/bucketStorage/TargetBucketName/TargetFileName .

```
## Points d'API

### Liste des Buckets
Méthode: GET
URL: /
Réponse: Liste des buckets en XML.


### Créer un Bucket
Méthode: PUT
URL: /{bucketName}
Corps: Configuration XML du bucket (optionnel).
Réponse: Confirmation de création du bucket.


### Supprimer un Bucket
Méthode: DELETE
URL: /{bucketName}
Réponse: Confirmation de suppression du bucket.


### Lister les Objets dans un Bucket
Méthode: GET
URL: /{bucketName}
Paramètres:
prefix (optionnel)
delimiter (optionnel)
max-keys (optionnel)
Réponse: Liste des objets en XML.

### Download un Objet dans un Bucket
Méthode: GET
URL: /{bucketName}/objects/{objectName}
Paramètres:
Réponse: Objet du bucket.


### Ajouter un Objet
Méthode: PUT
URL: /{bucketName}/{objectName}
FormData: Fichier à uploader.
Réponse: Métadonnées de l'objet ajouté.


### Supprimer un Objet
Méthode: DELETE
URL: /{bucketName}/{objectName}
Réponse: Confirmation de suppression de l'objet.


### Supprimer des Objets en Batch
Méthode: DELETE
URL: /{bucketName}
Corps: Liste des objets à supprimer en XML.
Réponse: Confirmation de suppression des objets en XML.

## Docker 
```sh 
docker-compose build
docker-compose up
```
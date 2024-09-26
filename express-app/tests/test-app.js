const chai = require('chai');
const should = chai.should();
const chaiHttp = require('chai-http');

chai.use(chaiHttp);

const app = require('../index.js'); // Assurez-vous que le chemin est correct

describe('GET /', () => {
    it('should respond with hello world', (done) => {
        chai.request(app)
            .get('/')
            .end((err, res) => {
                should.not.exist(err);
                res.status.should.equal(200);
                res.type.should.equal('text/plain');
                done();
            });
    });

    it('should respond with 404 for unknown route', (done) => {
        chai.request(app)
            .get('/unknown')
            .end((err, res) => {
                should.not.exist(err);
                res.status.should.equal(404);
                done();
            });
    });


    it('should return a 500 error for internal server error', (done) => {
        // Simuler une erreur interne pour un test
        chai.request(app)
            .get('/error') // Assurez-vous que vous avez une route qui renvoie une erreur
            .end((err, res) => {
                should.not.exist(err);
                res.status.should.equal(500);
                done();
            });
    });
});

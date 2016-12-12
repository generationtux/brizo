import { Injectable }       from '@angular/core';
import { tokenNotExpired }  from 'angular2-jwt';
import { Http }             from '@angular/http';

@Injectable()
export class AuthService {
    constructor(private http: Http) {}

    loggedIn() {
        return tokenNotExpired();
    }
}

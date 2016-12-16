import { Injectable }       from '@angular/core';
import { tokenNotExpired }  from 'angular2-jwt';
import { Headers, RequestOptions } from '@angular/http';

@Injectable()
export class AuthService {
  validateJwtToken() {
    return tokenNotExpired();
  }

  getToken() :string {
    return localStorage.getItem('id_token');
  }

  jwtRequestOptions() :RequestOptions {
    let headers = new Headers({'Authorization': 'Bearer ' + this.getToken()})
    return new RequestOptions({headers})
  }
}

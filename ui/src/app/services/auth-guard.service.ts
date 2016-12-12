import { Injectable }           from '@angular/core';
import { Router, CanActivate }  from '@angular/router';
import { AuthService }          from './auth.service';

@Injectable()
export class AuthGuard implements CanActivate {

    constructor(private auth: AuthService, private router: Router) {}

    canActivate() {
        if(this.auth.loggedIn()) {
            console.log('AUTH GUARD PASSED');
            return true;
        } else {
            console.log('BLOCKED BY AUTH GUARD');
            this.router.navigate(['/login']);
            return false;
        }
    }

}

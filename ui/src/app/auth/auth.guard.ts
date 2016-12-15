import { Injectable }           from '@angular/core';
import { Router, CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot }  from '@angular/router';
import { AuthService }          from './auth.service';
import {Observable} from 'rxjs/Rx';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private auth: AuthService, private router: Router) {}

  /*
  canActivate() {
    if(this.auth.validateJwtToken()) {
      return true;
    } else {
      this.router.navigate(['/login']);
      return false;
    }
  }
  */
  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | boolean {
    if (!this.auth.validateJwtToken()) {
      this.router.navigate(['/login']);
      return false;
    }
    return true;
  }
}

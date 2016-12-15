import { Injectable } from '@angular/core';
import { Router }     from '@angular/router';

@Injectable
export class AuthGuard implements CanActivate {
  constructor(private router: Router) {}

  //validateJwtToken()
}

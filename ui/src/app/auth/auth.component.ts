import { Component }              from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription }           from "rxjs/Rx";

@Component({
    selector: 'auth',
    templateUrl: './auth.component.html'
})

export class AuthComponent {
  private subscription: Subscription;
  private token: string;

  constructor(private router: Router, private route: ActivatedRoute) {
    this.subscription = route.queryParams.subscribe(
      (queryParam: any) => this.token = queryParam['token']
    );

    localStorage.setItem('id_token', this.token);
    this.router.navigate(['./app']);
  }

}

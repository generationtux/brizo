import { Component }              from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
    selector: 'auth',
    template: ''
})

export class AuthComponent {

  constructor(private router: Router, private route: ActivatedRoute) {
    route.queryParams.subscribe(
      (queryParam: any) => {
        localStorage.setItem('id_token', queryParam['token']);
        this.router.navigate(['./']);
      }
    );
  }

}

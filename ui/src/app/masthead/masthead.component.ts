import { Observable } from 'rxjs/Rx';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { Http, Headers, Response, RequestOptions } from '@angular/http';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';
import { ApplicationService } from '../resources/applications/application.service';
import { Router } from '@angular/router';

@Component({
    selector:       'masthead',
    styleUrls:      ['./masthead.css'],
    templateUrl:    './masthead.html',
    providers:      [ApplicationService],
})

export class MastheadComponent implements OnInit {
  public applications: any;
  public createAppForm: FormGroup;
  public submitted: boolean;
  public modalActions: EventEmitter<string>;

  constructor(private router: Router, private applicationService: ApplicationService, private _fb: FormBuilder) {
    this.modalActions = new EventEmitter<string>();
  }
  
  /*
  openModal() {
    this.modalActions.emit({action:"modal",params:['open']});
  }

  closeModal() {
    this.modalActions.emit({action:"modal",params:['close']});
  }
  */

  ngOnInit() {
    this.createAppForm = new FormGroup({
        name: new FormControl('', [<any>Validators.required])
    });
  }

  save() {
    return this.applicationService.createApplication(this.createAppForm.controls['name'].value).subscribe(
      err => console.error('There was an error: ' + err),
      () => (this.router.navigate([''])),
    )
  }
}

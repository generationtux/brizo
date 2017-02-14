import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { VersionService } from './version.service';
import { VersionDetailsComponent } from './details/version-details.component';
import { VersionCreateComponent } from './create/version-create.component';

const versionRoutes: Routes = [
  {
    path: 'environments/:environment-uuid',
    canActivate: [AuthGuard],
    children: [
      { path: 'versions/create', component: VersionCreateComponent, canActivate: [AuthGuard] },
      { path: 'versions/:version-uuid', component: VersionDetailsComponent, canActivate: [AuthGuard] },
    ]
  }
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(versionRoutes),
    SharedModule,
  ],
  declarations: [
    VersionDetailsComponent,
    VersionCreateComponent,
  ],
  providers: [
    VersionService,
  ],
  exports: [
    VersionDetailsComponent,
    VersionCreateComponent,
    RouterModule,
  ],
})

export class VersionModule {}

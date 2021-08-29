Write-Output "-----------backend begin-----------"
Start-Process -FilePath "wt" -ArgumentList "make", "runGo"
Start-Process -FilePath "wt" -ArgumentList  "pwsh", "./ui.ps1"

Set-Variable CONDA_PATH 'C:\Users\XuYipei\anaconda3'
Set-Variable PS1_PATH ($CONDA_PATH + '\shell\condabin\conda-hook.ps1')
Invoke-Expression $PS1_PATH
conda activate rcmdsys
python ./pkg/rcmdsys/server/main.py